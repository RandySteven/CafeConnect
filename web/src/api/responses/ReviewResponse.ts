interface ReviewsResponse {
    user: {
        id: number
        name: string
        profile_picture: string
    }
    score: number
    comment: string
    created_at: string
}

export interface GetReviewsResponse {
    cafe_id: number
    avg_score: number
    reviews: ReviewsResponse[]
}