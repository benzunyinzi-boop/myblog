<template>
  <section class="admin-page">
    <header class="admin-page-head">
      <div>
        <span class="section-chip">articles / manage</span>
        <h2>文章管理</h2>
        <p class="page-copy">支持分页、过滤、新建、编辑、发布与删除。</p>
      </div>
      <n-button type="primary" @click="handleCreate">+ 新建文章</n-button>
    </header>

    <section class="glass-card admin-filter-card">
      <n-space :size="12">
        <n-select
          v-model:value="filters.status"
          :options="statusOptions"
          placeholder="状态"
          clearable
          style="width: 140px"
          @update:value="handleFilterChange"
        />
        <n-select
          v-model:value="filters.category_id"
          :options="categoryOptions"
          placeholder="分类"
          clearable
          style="width: 160px"
          @update:value="handleFilterChange"
        />
        <n-input
          v-model:value="filters.keyword"
          placeholder="搜索标题/slug"
          clearable
          style="width: 200px"
          @keyup.enter="handleFilterChange"
        />
        <n-button @click="handleFilterChange">搜索</n-button>
      </n-space>
    </section>

    <section class="glass-card admin-table-card">
      <n-data-table
        :columns="columns"
        :data="articles"
        :loading="loading"
        :row-key="(row: ArticleSummary) => row.id"
        :pagination="pagination"
        @update:page="handlePageChange"
      />
    </section>
  </section>
</template>

<script setup lang="ts">
import { h, reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton,
  NDataTable,
  NInput,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  type DataTableColumns,
  type PaginationProps
} from 'naive-ui'
import { fetchCategories } from '../../api/meta'
import {
  adminListArticles,
  adminDeleteArticle,
  adminPublishArticle,
  adminUnpublishArticle,
  type ArticleSummary
} from '../../api/admin'
import { discrete } from '../../main'

const router = useRouter()

const articles = ref<ArticleSummary[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const filters = reactive({
  status: null as string | null,
  category_id: null as number | null,
  keyword: ''
})

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '已发布', value: 'published' }
]

const categoryOptions = ref<Array<{ label: string; value: number }>>([])

const pagination = computed<PaginationProps>(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    load()
  }
}))

const columns: DataTableColumns<ArticleSummary> = [
  { title: 'ID', key: 'id', width: 70 },
  {
    title: '标题',
    key: 'title',
    ellipsis: { tooltip: true },
    render: (row) => h('span', { style: 'font-weight:500;' }, row.title)
  },
  { title: 'Slug', key: 'slug', width: 160, ellipsis: true },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) =>
      h(
        NTag,
        { type: row.status === 'published' ? 'success' : 'default', size: 'small', bordered: false },
        { default: () => (row.status === 'published' ? '已发布' : '草稿') }
      )
  },
  { title: '浏览', key: 'view_count', width: 80 },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    render: (row) =>
      h(
        'div',
        { style: 'display:flex;gap:8px;flex-wrap:wrap;' },
        [
          h(
            NButton,
            { size: 'small', quaternary: true, onClick: () => handleEdit(row.id) },
            { default: () => '编辑' }
          ),
          row.status === 'draft'
            ? h(
                NButton,
                { size: 'small', quaternary: true, type: 'success', onClick: () => handlePublish(row.id) },
                { default: () => '发布' }
              )
            : h(
                NButton,
                { size: 'small', quaternary: true, type: 'warning', onClick: () => handleUnpublish(row.id) },
                { default: () => '下线' }
              ),
          h(
            NPopconfirm,
            { onPositiveClick: () => handleDelete(row.id) },
            {
              trigger: () =>
                h(NButton, { size: 'small', quaternary: true, type: 'error' }, { default: () => '删除' }),
              default: () => `确认删除文章「${row.title}」?`
            }
          )
        ]
      )
  }
]

async function load() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (filters.status) params.status = filters.status
    if (filters.category_id) params.category_id = filters.category_id
    if (filters.keyword?.trim()) params.keyword = filters.keyword.trim()

    const resp = await adminListArticles(params)
    articles.value = resp.data.items
    total.value = resp.data.total
  } catch (err) {
    discrete.message.error('加载文章失败')
    console.error(err)
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    const resp = await fetchCategories()
    categoryOptions.value = resp.data.items.map((c) => ({ label: c.name, value: c.id }))
  } catch (err) {
    console.error(err)
  }
}

function handleFilterChange() {
  currentPage.value = 1
  load()
}

function handlePageChange(page: number) {
  currentPage.value = page
  load()
}

function handleCreate() {
  router.push('/admin/articles/new')
}

function handleEdit(id: number) {
  router.push(`/admin/articles/${id}/edit`)
}

async function handlePublish(id: number) {
  try {
    const resp = await adminPublishArticle(id)
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '发布失败')
      return
    }
    discrete.message.success('已发布')
    await load()
  } catch (err: any) {
    const msg = err?.response?.data?.message || '发布失败'
    discrete.message.error(msg)
  }
}

async function handleUnpublish(id: number) {
  try {
    const resp = await adminUnpublishArticle(id)
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '下线失败')
      return
    }
    discrete.message.success('已下线')
    await load()
  } catch (err: any) {
    const msg = err?.response?.data?.message || '下线失败'
    discrete.message.error(msg)
  }
}

async function handleDelete(id: number) {
  try {
    const resp = await adminDeleteArticle(id)
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

onMounted(() => {
  loadCategories()
  load()
})
</script>
