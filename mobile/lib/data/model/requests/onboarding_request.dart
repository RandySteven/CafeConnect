import 'dart:typed_data';

class RegisterRequest {
  String? firstName;
  String? lastName;
  String? username;
  String? password;
  ByteBuffer? profilePicture;
  String? phoneNumber;
  String? dob;
  String? referralCode;
  String? address;
  Float64x2? longitude;
  Float64x2? latitude;
}

class LoginRequest {
  String? email;
  String? password;
}