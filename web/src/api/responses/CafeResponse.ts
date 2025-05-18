export interface ListCafeResponse {
    id : number
    name : string
    logo_url : string
    status : string
    open_hour : string
    close_hour : string
    address : string
}

export interface DetailCafeResponse {
    id: number
    name: string
    logo_url: string
    address:  {
        address: string
        latitude: number
        longitude: number
    }
    status: string
    photo_urls: string[]
    created_at: string
    updated_at: string
    deleted_at: string
}