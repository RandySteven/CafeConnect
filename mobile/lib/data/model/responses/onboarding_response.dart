class LoginResponse {
  String? accessToken;
  String? refreshToken;
  String? loginTime;

  LoginResponse({
    this.accessToken, this.refreshToken, this.loginTime
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      accessToken: json['access_token'],
      refreshToken: json['refresh_token'],
      loginTime: json['login_time']
    );
  }
}

class RegisterResponse {
  String? id;
  String? email;
  String? registerTime;

  RegisterResponse({
    this.id, this.email, this.registerTime
  });

  factory RegisterResponse.fromJson(Map<String, dynamic> json) {
    return RegisterResponse(
      id: json['id'],
      email: json['email'],
      registerTime: json['register_time']
    );
  }
}