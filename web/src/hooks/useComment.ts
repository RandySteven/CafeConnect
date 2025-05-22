import {GetReviewsResponse} from "@/api/responses/ReviewResponse";
import {useEffect, useState} from "react";
import {POST} from "@/api/api";
import {GetReviewCafeRequest} from "@/api/requests/ReviewRequest";
import {GET_CAFE_REVIEWS} from "@/api/endpoint";

export const useGetReviews = (id : number) : GetReviewsResponse => {
    const [reviewsResponse, setReviewsResponse] = useState<GetReviewsResponse>({
        cafe_id: 0,
        avg_score: 0,
        reviews: []
    })

    let request : GetReviewCafeRequest = {
        cafe_id: Number(id)
    }


    useEffect(() => {
        const fetchReviews = async () => {
            try {
                const result = await POST(GET_CAFE_REVIEWS, false, request)
                if(result.data.review != null) {
                    setReviewsResponse(result.data.review)
                }
            }catch (error) {
                console.log(error)
            }
        }

        fetchReviews()
    }, []);

    return reviewsResponse
}