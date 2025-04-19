import 'package:mobile/data/model/requests/login_request.dart';
import 'package:mobile/data/model/requests/register_request.dart';
import 'package:mobile/data/model/responses/login_response.dart';
import 'package:mobile/data/model/responses/register_response.dart';

abstract class OnboardingApiRepository {
  Future<RegisterResponse> register(RegisterRequest request);
  Future<LoginResponse> login(LoginRequest request);
}