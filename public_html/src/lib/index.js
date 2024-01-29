// place files you want to import through the `$lib` alias in this folder.

Object.defineProperty(Number.prototype, 'humanReadableSize', {
    value: function (a, b, c, d) {
        return (a = a ? [1e3, 'k', 'B'] : [1024, 'K', 'iB'], b = Math, c = b.log,
                d = c(this) / c(a[0]) | 0, this / b.pow(a[0], d)).toFixed(2)
            + ' ' + (d ? (a[1] + 'MGTPEZY')[--d] + a[2] : 'Bytes');
    }, writable: false, enumerable: false
});

Object.defineProperty(String.prototype, 'Extension', {
    value: function () {
        return this.split('.').pop();
    }, writable: false
});

Object.defineProperty(String.prototype, 'IsImage', {
    value: function () {
        return ['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/bmp'].includes(this);
    }, writable: false
});

Object.defineProperty(String.prototype, 'IsVideo', {
    value: function () {
        return ['video/mp4', 'video/webm', 'video/ogg'].includes(this);
    }, writable: false
});

Object.defineProperty(String.prototype, 'Hash256', {
    value: async function () {
        const utf8 = new TextEncoder().encode(this);
        return await crypto.subtle.digest('SHA-256', utf8).then((hashBuffer) => {
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            const hashHex = hashArray
                .map((bytes) => bytes.toString(16).padStart(2, '0'))
                .join('');
            return hashHex;
        });
    }, writable: false
})