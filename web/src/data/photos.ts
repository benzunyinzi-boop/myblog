export type PhotoTag = 'street' | 'landscape' | 'portrait' | 'mono'

export type Photo = {
  id: number
  src: string
  title: string
  exif: string
  tag: PhotoTag
  takenAt: string
}

export const photoFilters: Array<{ label: string; value: PhotoTag | 'all' }> = [
  { label: '全部', value: 'all' },
  { label: '街头', value: 'street' },
  { label: '风景', value: 'landscape' },
  { label: '人像', value: 'portrait' },
  { label: '黑白', value: 'mono' }
]

export const photos: Photo[] = [
  {
    id: 1,
    src: 'https://images.unsplash.com/photo-1519501025264-65ba15a82390?w=1200&q=80',
    title: 'Midnight Station',
    exif: 'f/1.8 · 1/60 · ISO 1600',
    tag: 'street',
    takenAt: '2024-11'
  },
  {
    id: 2,
    src: 'https://images.unsplash.com/photo-1470252649378-9c29740c9fa8?w=1200&q=80',
    title: 'Mountain Silence',
    exif: 'f/8 · 1/250 · ISO 100',
    tag: 'landscape',
    takenAt: '2024-08'
  },
  {
    id: 3,
    src: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=1200&q=80',
    title: 'Portrait No.7',
    exif: 'f/1.4 · 1/200 · ISO 400',
    tag: 'portrait',
    takenAt: '2024-05'
  },
  {
    id: 4,
    src: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=1200&q=80',
    title: 'After the Rain',
    exif: 'f/2.8 · 1/500 · ISO 200',
    tag: 'landscape',
    takenAt: '2024-07'
  },
  {
    id: 5,
    src: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1200&q=80',
    title: 'Fog over Ridge',
    exif: 'f/5.6 · 1/400 · ISO 200',
    tag: 'mono',
    takenAt: '2024-02'
  },
  {
    id: 6,
    src: 'https://images.unsplash.com/photo-1519125323398-675f0ddb6308?w=1200&q=80',
    title: 'Neon Alley',
    exif: 'f/1.8 · 1/125 · ISO 800',
    tag: 'street',
    takenAt: '2024-10'
  },
  {
    id: 7,
    src: 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=1200&q=80',
    title: 'Forest Path',
    exif: 'f/4 · 1/160 · ISO 400',
    tag: 'landscape',
    takenAt: '2024-06'
  },
  {
    id: 8,
    src: 'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=1200&q=80',
    title: 'Window Light',
    exif: 'f/2 · 1/320 · ISO 200',
    tag: 'portrait',
    takenAt: '2024-03'
  }
]
