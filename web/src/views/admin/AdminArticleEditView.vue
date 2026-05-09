<template>
  <section class="admin-page">
    <header class="admin-page-head">
      <div>
        <span class="section-chip">article / editor</span>
        <h2>{{ isNew ? '新建文章' : '编辑文章' }}</h2>
        <p class="page-copy">支持 Markdown 编辑、图片上传、分类标签选择。</p>
      </div>
      <n-space :size="12">
        <n-button @click="handleBack">返回列表</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
      </n-space>
    </header>

    <section class="glass-card admin-form-card">
      <n-form :model="form" label-placement="top">
        <div class="form-grid">
          <n-form-item label="标题" required>
            <n-input v-model:value="form.title" placeholder="文章标题" />
          </n-form-item>
          <n-form-item label="Slug" required>
            <n-input v-model:value="form.slug" placeholder="url-friendly-slug" />
          </n-form-item>
        </div>

        <n-form-item label="摘要">
          <n-input v-model:value="form.summary" type="textarea" :rows="2" placeholder="一句话摘要(可选)" />
        </n-form-item>

        <div class="form-grid">
          <n-form-item label="分类">
            <n-select
              v-model:value="form.category_id"
              :options="categoryOptions"
              placeholder="选择分类"
              clearable
            />
          </n-form-item>
          <n-form-item label="标签">
            <n-select
              v-model:value="form.tag_ids"
              :options="tagOptions"
              placeholder="选择标签"
              multiple
              clearable
            />
          </n-form-item>
        </div>

        <n-form-item label="封面图">
          <n-input v-model:value="form.cover_image" placeholder="/uploads/... 或外链">
            <template #suffix>
              <n-upload
                :custom-request="handleUpload"
                :show-file-list="false"
                accept="image/*"
              >
                <n-button size="small" quaternary :loading="uploading">上传</n-button>
              </n-upload>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="正文 (Markdown)" required>
          <n-input
            v-model:value="form.content"
            type="textarea"
            :rows="18"
            placeholder="支持 Markdown 语法"
          />
        </n-form-item>

        <n-form-item label="状态">
          <n-radio-group v-model:value="form.status">
            <n-radio value="draft">草稿</n-radio>
            <n-radio value="published">已发布</n-radio>
          </n-radio-group>
        </n-form-item>
      </n-form>
    </section>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NRadio,
  NRadioGroup,
  NSelect,
  NSpace,
  NUpload,
  type UploadCustomRequestOptions
} from 'naive-ui'
import { fetchCategories, fetchTags } from '../../api/meta'
import {
  adminGetArticle,
  adminCreateArticle,
  adminUpdateArticle,
  adminUploadFile,
  type ArticlePayload
} from '../../api/admin'
import { discrete } from '../../main'

const route = useRoute()
const router = useRouter()

const articleId = computed(() => {
  const id = route.params.id
  return id ? Number(id) : null
})

const isNew = computed(() => !articleId.value)

const saving = ref(false)
const uploading = ref(false)

const form = reactive<ArticlePayload>({
  title: '',
  slug: '',
  summary: '',
  content: '',
  cover_image: '',
  category_id: undefined,
  tag_ids: [],
  status: 'draft'
})

const categoryOptions = ref<Array<{ label: string; value: number }>>([])
const tagOptions = ref<Array<{ label: string; value: number }>>([])

async function load() {
  try {
    const [catResp, tagResp] = await Promise.all([fetchCategories(), fetchTags()])
    categoryOptions.value = catResp.data.items.map((c) => ({ label: c.name, value: c.id }))
    tagOptions.value = tagResp.data.items.map((t) => ({ label: t.name, value: t.id }))

    if (!isNew.value && articleId.value) {
      const resp = await adminGetArticle(articleId.value)
      const article = resp.data
      form.title = article.title
      form.slug = article.slug
      form.summary = article.summary || ''
      form.content = article.content
      form.cover_image = article.cover_image || ''
      form.category_id = article.category_id || undefined
      form.tag_ids = article.tags.map((t) => t.id)
      form.status = article.status as 'draft' | 'published'
    }
  } catch (err) {
    discrete.message.error('加载数据失败')
    console.error(err)
  }
}

async function handleSave() {
  if (!form.title.trim() || !form.slug.trim() || !form.content.trim()) {
    discrete.message.warning('标题、Slug、正文不能为空')
    return
  }
  saving.value = true
  try {
    const payload: ArticlePayload = {
      title: form.title.trim(),
      slug: form.slug.trim(),
      summary: form.summary?.trim(),
      content: form.content.trim(),
      cover_image: form.cover_image?.trim(),
      category_id: form.category_id,
      tag_ids: form.tag_ids,
      status: form.status
    }
    const resp = isNew.value
      ? await adminCreateArticle(payload)
      : await adminUpdateArticle(articleId.value!, payload)
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '保存失败')
      return
    }
    discrete.message.success(isNew.value ? '创建成功' : '已保存')
    router.push('/admin/articles')
  } catch (err: any) {
    const msg = err?.response?.data?.message || '保存失败'
    discrete.message.error(msg)
  } finally {
    saving.value = false
  }
}

async function handleUpload(options: UploadCustomRequestOptions) {
  const { file } = options
  if (!file.file) return
  uploading.value = true
  try {
    const resp = await adminUploadFile(file.file as File)
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '上传失败')
      return
    }
    form.cover_image = resp.data.url
    discrete.message.success('上传成功')
  } catch (err: any) {
    const msg = err?.response?.data?.message || '上传失败'
    discrete.message.error(msg)
  } finally {
    uploading.value = false
  }
}

function handleBack() {
  router.push('/admin/articles')
}

onMounted(load)
</script>
