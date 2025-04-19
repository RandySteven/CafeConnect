import 'package:flutter/material.dart';

Widget InputFormWidget({
  required String label,
  required TextEditingController controller,
  bool obscureText = false,
}) {
  return TextField(
    controller: controller,
    obscureText: obscureText,
    decoration: InputDecoration(
      labelText: label,
      border: OutlineInputBorder()
    ),
  );
}