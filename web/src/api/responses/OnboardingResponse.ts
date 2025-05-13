import {Property} from "csstype";
import Float = Property.Float;

export interface LoginResponse {
    access_token : String
    refresh_token : String
    login_time : String
}

export interface RegisterResponse {
    id : String
    email : String
    register_time : String
}

interface OnboardAddressUser {
    id : Number
    address : String
    longitude : Float
    latitude : Float
    is_default : Boolean
}

export interface OnboardUserResponse {
    id : Number
    name : String
    username : String
    email : String
    profile_picture : String
    point : Number
    addresses : OnboardAddressUser[]
    created_at : String
    updated_at : String
    deleted_at : String
}