jQuery('#qrcodeCanvas').qrcode({
    text: "http://jetienne.com/"
});

function getExtention(fname) {
    return fname.slice((fname.lastIndexOf(".") - 1 >>> 0) + 2);
}

function pathJoin(parts, sep) {
    var separator = sep || '/';
    var replace = new RegExp(separator + '{1,}', 'g');
    return parts.join(separator).replace(replace, separator);
}

var vm = new Vue({
    el: "#app",
    data: {
        message: "Hello vue.js",
        breadcrumb: [],
        showHidden: false,
        files: [{
            name: "loading ...",
            path: "",
            size: "...",
            type: "dir",
        }]
    },
    computed: {
        computedFiles: function() {
            var that = this;
            return this.files.filter(function(f) {
                if (that.showHidden && f.name.slice(0, 1) === '.') {
                    return false;
                }
                return true;
            });
        },
    },
    methods: {
        toggleHidden: function() {
            this.showHidden = !this.showHidden;
        },
        genQrcode: function(text) {
            var urlPath = location.protocol + "//" + location.host + location.pathname + text;
            $("#qrcode-title").html(text);
            $("#qrcode-link").attr("href", urlPath);
            $('#qrcodeCanvas').empty().qrcode({
                text: urlPath
            });
            $("#qrcode-modal").modal("show");
        },
        genFileClass: function(f) {
            if (f.type == "dir") {
                if (f.name == '.git') {
                    return 'fa-git-square';
                }
                return "fa-folder-open";
            }
            var ext = getExtention(f.name);
            switch (ext) {
                case "go":
                case "py":
                case "js":
                case "java":
                case "c":
                case "cpp":
                case "h":
                    return "fa-file-code-o";
                case "pdf":
                    return "fa-file-pdf-o";
                case "zip":
                    return "fa-file-zip-o";
                case "mp3":
                case "wav":
                    return "fa-file-audio-o";
                case "jpg":
                case "png":
                case "gif":
                case "jpeg":
                case "tiff":
                    return "fa-file-picture-o";
            }
            return "fa-file-text-o"
        },
        clickFileOrDir: function(f, e) {
            if (f.type == "file") {
                return true;
            }
            var reqPath = pathJoin([location.pathname, f.name]);
            loadDirectory(reqPath);
            e.preventDefault()
        },
        changePath: function(reqPath, e) {
            loadDirectory(reqPath);
            e.preventDefault()
        },
        updateBreadcrumb: function() {
            var pathname = decodeURI(location.pathname || "/");
            var parts = pathname.split('/');
            this.breadcrumb = [];
            if (pathname == "/") {
                return this.breadcrumb;
            }
            var i = 2;
            for (; i <= parts.length; i += 1) {
                var name = parts[i - 1];
                var path = parts.slice(0, i).join('/');
                this.breadcrumb.push({
                    name: name + (i == parts.length ? ' /' : ''),
                    path: path
                })
            }
            return this.breadcrumb;
        }
    }
})

window.onpopstate = function(event) {
    var pathname = decodeURI(location.pathname)
    loadFileList()
}

function loadDirectory(reqPath) {
    window.history.pushState({}, "", reqPath);
    loadFileList()
}

function loadFileList() {
    $.getJSON("/-/json" + location.pathname, function(res) {
        // console.log(res)
        res.sort(function(a, b) {
            var obj2n = function(v) {
                return v.type == "dir" ? 0 : 1;
            }
            return obj2n(a) - obj2n(b);
        })
        vm.files = res;
    })
    vm.updateBreadcrumb()
}
loadFileList() // For page first loading
