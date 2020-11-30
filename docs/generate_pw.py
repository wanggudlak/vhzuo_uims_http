# -*- coding: utf-8 -*-

import hmac
import random
import hashlib
from struct import Struct

_pack_int = Struct('>I').pack
izip = zip

from binascii import b2a_hex as _b2a_hex


def b2a_hex(s):
    return _b2a_hex(s).decode('us-ascii')


def binxor(a, b):
    return bytes([x ^ y for (x, y) in zip(a, b)])


def pbkdf2_hex(data, salt, iterations=1000, keylen=24, hashfunc=None):
    passwd_bytes = pbkdf2_bin(data, salt, iterations, keylen, hashfunc)
    # print(passwd_bytes)
    return b2a_hex(passwd_bytes)


def pbkdf2_bin(data, salt, iterations=1000, keylen=24, hashfunc=None):
    """
    Returns a binary digest for the PBKDF2 hash algorithm of `data`
    with the given `salt`.  It iterates `iterations` time and produces a
    key of `keylen` bytes.  By default SHA-1 is used as hash function,
    a different hashlib `hashfunc` can be provided.
    """
    hashfunc = hashfunc or hashlib.sha1
    # 将原密码和盐值转换成bytes类型处理
    translateData = data.encode("utf-8") if isinstance(data, str) else data
    translateSolt = salt.encode("utf-8") if isinstance(salt, str) else salt
    mac = hmac.new(translateData, None, hashfunc)
    def _pseudorandom(x, mac=mac):
        h = mac.copy()
        h.update(x)
        return [n for n in bytearray(h.digest())]
    buf = []
    for block in range(1, -(-keylen // mac.digest_size) + 1):
        rv = u = _pseudorandom(translateSolt + _pack_int(block))
        for i in range(iterations - 1):
            b1 = bytes()
            for c in u:
                b1 += bytes([c])
            u = _pseudorandom(b1)
            rv = binxor(rv, u)
        buf.extend(rv)
    # 之前的内置方法chr会直接将数字转换成单个字符,且字符的编码方式有所不同
    # 转换成Bytes类型，拼接完后直接处理
    b2 = bytes()
    for c in buf:
        b2 += bytes([c])
    return b2[:keylen]


def rand_str(num=10):
    return "".join(random.sample("ABCDEFGHJKLMNPQRSTUVWXY23456789ABCDEFGHJKLMNPQRSTUVWXY23456789abcdefghjkmnpqrstuvwxy23456789abcdefghjkmnpqrstuvwxy23456789", num))



def generate_password(origin_password):
    salt = rand_str()
    pw = pbkdf2_hex(origin_password, salt, 10000)
    return pw, salt


def check_password(origin_password, encrypt_password, salt):
    if len(salt) > 15:
        sign = hashlib.md5("%s%s" % (salt, origin_password)).hexdigest()
        return sign == encrypt_password
    return pbkdf2_hex(origin_password, salt, 10000) == encrypt_password



if __name__ == '__main__':
    # 生成密码及盐值
    pw, salt = generate_password("123456")
    print("%s, %s" % (pw, salt))


    # 验证密码
    result = check_password("123456", "261bceeefd0d72f1b5cf98993798f6056b8ebfa26265dd26", "9AbQYARjuk")
    print(result)
