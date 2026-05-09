<template>
  <section class="admin-page">
    <header class="admin-page-head">
      <div>
        <span class="section-chip">categories / manage</span>
        <h2>分类管理</h2>
        <p class="page-copy">对应 /admin/categories,支持创建、编辑、删除。sort_order 越小越靠前。</p>
      </div>
      <n-button type="primary" @click="openCreate">+ 新建分类</n-button>
    </header>

    <section class="glass-card admin-table-card">
      <n-data-table
        :columns="columns"
        :data="categories"
        :loading="loading"
        :row-key="(row: Category) => row.id"
      />
    </section>

    <n-modal
      v-model:show="dialogOpen"
      preset="card"
      :title="editing ? '编辑分类' : '新建分类'"
      class="admin-modal"
      :style="{ width: '460px' }"
      @close="handleClose"
    >
      <n-form :model="form" label-placement="top">
        <n-form-item label="名称" required>
          <n-input v-model:value="form.name" placeholder="例如:Golang" />
        </n-form-item>
        <n-form-item label="Slug" required>
          <n-input v-model:value="form.slug" placeholder="小写字母-数字,例如 golang" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="form.description" type="textarea" :rows="2" placeholder="一句话描述(可选)" />
        </n-form-item>
        <n-form-item label="排序">
          <n-input-number v-model:value="form.sort_order" :min="-1000" :max="1000" />
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
  NInputNumber,
  NModal,
  NPopconfirm,
  type DataTableColumns
} from 'naive-ui'
import { fetchCategories, type Category } from '../../api/meta'
import { adminCreateCategory, adminUpdateCategory, adminDeleteCategory } from '../../api/admin'
import { discrete } from '../../main'

const categories = ref<Category[]>([])
const loading = ref(false)

const dialogOpen = ref(false)
const submitting = ref(false)
const editing = ref<Category | null>(null)
const form = reactive({ name: '', slug: '', description: '', sort_order: 0 })

const columns: DataTableColumns<Category> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', key: 'name' },
  { title: 'Slug', key: 'slug' },
  { title: '描述', key: 'description', ellipsis: true },
  { title: '排序', key: 'sort_order', width: 90 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: (row) =>
      h(
        'div',
        { style: 'display:flex;gap:8px;' },
        [
          h(
            NButton,
            { size: 'small', quaternary: true, onClick: () => openEdit(row) },
            { default: () => '编辑' }
          ),
          h(
            NPopconfirm,
            { onPositiveClick: () => handleDelete(row.id) },
            {
              trigger: () =>
                h(NButton, { size: 'small', quaternary: true, type: 'error' }, { default: () => '删除' }),
              default: () => `确认删除分类「${row.name}」?`
            }
          )
        ]
      )
  }
]

async function load() {
  loading.value = true
  try {
    const resp = await fetchCategories()
    categories.value = resp.data.items
  } catch (err) {
    discrete.message.error('加载分类失败')
    console.error(err)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  form.name = ''
  form.slug = ''
  form.description = ''
  form.sort_order = 0
  dialogOpen.value = true
}

function openEdit(row: Category) {
  editing.value = row
  form.name = row.name
  form.slug = row.slug
  form.description = row.description || ''
  form.sort_order = row.sort_order || 0
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
    const payload = {
      name: form.name.trim(),
      slug: form.slug.trim(),
      description: form.description.trim(),
      sort_order: form.sort_order
    }
    const resp = editing.value
      ? await adminUpdateCategory(editing.value.id, payload)
      : await adminCreateCategory(payload)
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '保存失败')
      return
    }
    discrete.message.success(editing.value ? '已更新' : '创建成功')
    dialogOpen.value = false
    await load()
  } catch (err: any) {
    const msg = err?.response?.data?.message || '保存失败'
    discrete.message.error(msg)
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: number) {
  try {
    const resp = await adminDeleteCategory(id)
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
