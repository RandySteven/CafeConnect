import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:mobile/constant/endpoint.dart';

import 'package:mobile/data/model/requests/onboarding_request.dart';
import 'package:mobile/data/model/responses/onboarding_response.dart';

abstract class OnboardingApiRepository {
  Future<RegisterResponse> register(RegisterRequest request);
  Future<LoginResponse> login(LoginRequest request);
}

class OnboardingApiRepositoryImpl implements OnboardingApiRepository {

  @override
  Future<LoginResponse> login(LoginRequest request) async {
    var response = await http.post((BASE_URL+LOGIN_API) as Uri, headers: {
      "Content-Type": "application/json"
    }, body: jsonEncode({
      "email": request.email,
      "password": request.password
    }));
    if(response.statusCode == 200) {
      final Map<String, dynamic> result = jsonDecode(response.body);
      return LoginResponse.fromJson(result);
    }else {
      throw Exception('Failed to login: ${response.statusCode} ${response.body}');
    }
  }

  @override
  Future<RegisterResponse> register(RegisterRequest request) async {
    var response = await http.post(
      (BASE_URL + REGISTER_API) as Uri,
      headers: {
        "Content-Type": "multipart/form-data"
      },
      body: null
    );
    if(response.statusCode == 201) {
      final Map<String, dynamic> result = jsonDecode(response.body);
      return RegisterResponse.fromJson(result);
    }else {
      throw Exception('Failed to register: ${response.statusCode} ${response.body}');
    }
  }
  
}