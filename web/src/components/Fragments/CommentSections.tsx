"use client";

import {Fragment} from "react";
import {Box} from "@mui/material";
import {useGetReviews} from "@/hooks/useComment";
import {CommentCard} from "@/components/Elements/Card";

export const CommentSections = (prop : {
    id: number
}) => {
    const reviews = useGetReviews(prop.id)
    return <Fragment>
        <Box>
            {
                reviews.reviews.map((review, index) => (
                    <CommentCard
                        key={index}
                        avatar={review.user.profile_picture}
                        name={review.user.name}
                        score={review.score}
                        comment={review.comment}
                        timestamp={review.created_at} />
                ))
            }
        </Box>
    </Fragment>
}