import service from '@/apis/client/service'

export interface MenusRequest {
  userId: string
}

export interface menuItemData {
  name: string
  path: string
  component?: string
  meta?: { [key: string]: any } | undefined  // 修改这里，使用对象而不是Map
  children?: menuItemData[]
}

export const fetchMenus = (req: MenusRequest) => {
  return service.get<menuItemData[]>('/api/system/menus', req, {
    silent: false
  })
}