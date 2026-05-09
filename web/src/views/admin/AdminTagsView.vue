<template>
  <section class="admin-page">
    <header class="admin-page-head">
      <div>
        <span class="section-chip">tags / manage</span>
        <h2>标签管理</h2>
        <p class="page-copy">对应 /admin/tags,支持创建和删除。标签与文章多对多,删除时请确认没有关联文章。</p>
      </div>
      <n-button type="primary" @click="openCreate">+ 新建标签</n-button>
    </header>

    <section class="glass-card admin-table-card">
      <n-data-table
        :columns="columns"
        :data="tags"
        :loading="loading"
        :row-key="(row: Tag) => row.id"
      />
    </section>

    <n-modal
      v-model:show="dialogOpen"
      preset="card"
      title="新建标签"
      class="admin-modal"
      :style="{ width: '420px' }"
      @close="handleClose"
    >
      <n-form :model="form" label-placement="top">
        <n-form-item label="名称" required>
          <n-input v-model:value="form.name" placeholder="例如:并发" />
        </n-form-item>
        <n-form-item label="Slug" required>
          <n-input v-model:value="form.slug" placeholder="小写字母-数字,例如 concurrency" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="modal-footer">
          <n-button @click="handleClose">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">保存</n-button>
        </div>
      </template>
    </n-modal>
  </section>
</template>

<script setup lang="ts">
import { h, reactive, ref, onMounted } from 'vue'
import {
  NButton,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NPopconfirm,
  NTag,
  type DataTableColumns
} from 'naive-ui'
import { fetchTags, type Tag } from '../../api/meta'
import { adminCreateTag, adminDeleteTag } from '../../api/admin'
import { discrete } from '../../main'

const tags = ref<Tag[]>([])
const loading = ref(false)

const dialogOpen = ref(false)
const submitting = ref(false)
const form = reactive({ name: '', slug: '' })

const columns: DataTableColumns<Tag> = [
  { title: 'ID', key: 'id', width: 80 },
  {
    title: '名称',
    key: 'name',
    render: (row) => h(NTag, { type: 'info', bordered: false }, { default: () => row.name })
  },
  { title: 'Slug', key: 'slug' },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row) =>
      h(
        NPopconfirm,
        {
          onPositiveClick: () => handleDelete(row.id)
        },
        {
          trigger: () =>
            h(NButton, { size: 'small', quaternary: true, type: 'error' }, { default: () => '删除' }),
          default: () => `确认删除标签「${row.name}」?`
        }
      )
  }
]

async function load() {
  loading.value = true
  try {
    const resp = await fetchTags()
    tags.value = resp.data.items
  } catch (err) {
    discrete.message.error('加载标签失败')
    console.error(err)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.name = ''
  form.slug = ''
  dialogOpen.value = true
}

function handleClose() {
  dialogOpen.value = false
}

async function handleSubmit() {
  if (!form.name.trim() || !form.slug.trim()) {
    discrete.message.warning('请填写名称和 slug')
    return
  }
  submitting.value = true
  try {
    const resp = await adminCreateTag({ name: form.name.trim(), slug: form.slug.trim() })
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '创建失败')
      return
    }
    discrete.message.success('创建成功')
    dialogOpen.value = false
    await load()
  } catch (err: any) {
    const msg = err?.response?.data?.message || '创建失败'
    discrete.message.error(msg)
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: number) {
  try {
    const resp = await adminDeleteTag(id)
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '删除失败')
      return
    }
    discrete.message.success('已删除')
    await load()
  } catch (err: any) {
    const msg = err?.response?.data?.message || '删除失败'
    discrete.message.error(msg)
  }
}

onMounted(load)
</script>
