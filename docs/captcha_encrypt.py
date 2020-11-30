# -*- coding: utf-8 -*-

import math

ALPHABET= "abcdefghijklmnopqrstuvwxyz"

# 图片坐标加密
def encrypt_captcha(value):
	ascii_list = list(map(ord, value))
	return [{ALPHABET[i]:pow(v,2)} for i, v in  enumerate(ascii_list)]


# 图片坐标解密
def decrypt_captcha(value):
	# 前端加密 x = x*2 +y*2 y =x*2
	aa = [int(math.sqrt(v[ALPHABET[i]])) for i, v in enumerate(value)]
	print("aa=", aa)

	maa = map(chr, aa)

	char= "".join(maa)
	print("char=", char)

	i = char.find(":")
	a, b = int(math.sqrt(int(char[i+1:]))), int(math.sqrt(int(int(char[0:i]))-(int(char[i+1:]))))
	return a, b


if __name__ == '__main__':
	# 加密坐标值
	# encr = encrypt_captcha("100:50")
	# print("%s" % encr)


	# 解密坐标值
	result = decrypt_captcha([
		{"a":2916},
		{"b":2401},
		{"c":2500},
		{"d":2304},
		{"e":2401},
		{"f":3364},
		{"g":2809},
		{"h":2304},
		{"i":2401},
		{"j":3025},
		{"k":2916},])
	print(result)