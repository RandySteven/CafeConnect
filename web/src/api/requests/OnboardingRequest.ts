import {Property} from "csstype";
import Float = Property.Float;

export interface LoginRequest {
    email : String
    password : String
}

export interface RegisterRequest {
    first_name : String
    last_name : String
    username : String
    email :  String
    password : String
    profile_picture : File
    phone_number : String
    dob : String
    referral_code : String
    address : String
    longitude : Float
    latitude : Float
}