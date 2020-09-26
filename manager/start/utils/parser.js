function stringToUint8Array(string, options = {stream: false}) {
    if (options.stream) {
        throw new Error(`Failed to encode: the 'stream' option is unsupported.`);
    }

    var pos = 0;
    const len = string.length;
    const out = [];

    var at = 0;  // output position
    var tlen = Math.max(32, len + (len >> 1) + 7);  // 1.5x size
    var target = new Uint8Array((tlen >> 3) << 3);  // ... but at 8 byte offset

    while (pos < len) {
        var value = string.charCodeAt(pos++);
        if (value >= 0xd800 && value <= 0xdbff) {
            // high surrogate
            if (pos < len) {
                const extra = string.charCodeAt(pos);
                if ((extra & 0xfc00) === 0xdc00) {
                    ++pos;
                    value = ((value & 0x3ff) << 10) + (extra & 0x3ff) + 0x10000;
                }
            }
            if (value >= 0xd800 && value <= 0xdbff) {
                continue;  // drop lone surrogate
            }
        }

        // expand the buffer if we couldn't write 4 bytes
        if (at + 4 > target.length) {
            tlen += 8;  // minimum extra
            tlen *= (1.0 + (pos / string.length) * 2);  // take 2x the remaining
            tlen = (tlen >> 3) << 3;  // 8 byte offset

            const update = new Uint8Array(tlen);
            update.set(target);
            target = update;
        }

        if ((value & 0xffffff80) === 0) {  // 1-byte
            target[at++] = value;  // ASCII
            continue;
        } else if ((value & 0xfffff800) === 0) {  // 2-byte
            target[at++] = ((value >> 6) & 0x1f) | 0xc0;
        } else if ((value & 0xffff0000) === 0) {  // 3-byte
            target[at++] = ((value >> 12) & 0x0f) | 0xe0;
            target[at++] = ((value >> 6) & 0x3f) | 0x80;
        } else if ((value & 0xffe00000) === 0) {  // 4-byte
            target[at++] = ((value >> 18) & 0x07) | 0xf0;
            target[at++] = ((value >> 12) & 0x3f) | 0x80;
            target[at++] = ((value >> 6) & 0x3f) | 0x80;
        } else {
            // FIXME: do we care
            continue;
        }

        target[at++] = (value & 0x3f) | 0x80;
    }

    return target.slice(0, at);
}

// 字节到文本
function Utf8ArrayToStr(array) {
    var out, i, len, c;
    var char2, char3;

    out = "";
    len = array.length;
    i = 0;
    while (i < len) {
        c = array[i++];
        switch (c >> 4) {
            case 0:
            case 1:
            case 2:
            case 3:
            case 4:
            case 5:
            case 6:
            case 7:
                // 0xxxxxxx
                out += String.fromCharCode(c);
                break;
            case 12:
            case 13:
                // 110x xxxx   10xx xxxx
                char2 = array[i++];
                out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));
                break;
            case 14:
                // 1110 xxxx  10xx xxxx  10xx xxxx
                char2 = array[i++];
                char3 = array[i++];
                out += String.fromCharCode(((c & 0x0F) << 12) |
                    ((char2 & 0x3F) << 6) |
                    ((char3 & 0x3F) << 0));
                break;
        }
    }

    return out;
}

// 解析服务端的响应
function blobToUint8Array(b) {
    var uri = URL.createObjectURL(b),
        xhr = new XMLHttpRequest()

    xhr.open('GET', uri, false);
    xhr.send();

    URL.revokeObjectURL(uri);

    return base64XORToUint8Array(xhr.response);
}

// 解析成对象
function deserializeObject(data) {
    try {
        return proto.protobuf.Response.deserializeBinary(blobToUint8Array(data));
    } catch (e) {
        return Utf8ArrayToStr(blobToUint8Array(data))
    }
}

// Base64 XOR 转 array
function base64XORToUint8Array(base64String) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
        .replace(/\-/g, '+')
        .replace(/_/g, '/');

    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);

    for (var i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i) ^ 32;
    }
    return outputArray;
}

function XORData(data) {
    const outputArray = new Uint8Array(data.length);
    for (var i = 0; i < data.length; ++i) {
        outputArray[i] = data[i] ^ 32;
    }
    return outputArray
}

function getReqId(data) {
    try {
        return data.getId()
    } catch (e) {
        return ""
    }
}

// 计算签名
function calcSign(data) {
    return stringToUint8Array(md5(data));
}

function getRequest(path, type, location, data) {
    var t = new proto.protobuf.Request();
    var timestamp = (new Date()).getTime();
    t.setPath(path);
    t.setType(type);
    t.setLocation(location);
    // 开始计算签名
    data.salt = timestamp.toString();
    var saltArray = stringToUint8Array(timestamp.toString());
    var dataArray = stringToUint8Array(JSON.stringify(data));
    var signArray = new Uint8Array(saltArray.length + dataArray.length);
    signArray.set(saltArray);
    signArray.set(dataArray, saltArray.length);
    t.setSign(calcSign(signArray));
    // 计算签名结束
    t.setData(dataArray);
    return t
}

//
// export {
//     stringToUint8Array,
//     Utf8ArrayToStr,
//     blobToUint8Array,
//     base64XORToUint8Array,
//     XORData,
//     calcSign
// }