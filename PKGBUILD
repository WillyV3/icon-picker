# Maintainer: WillyV3 <your-email@example.com>
pkgname=icon-picker
pkgver=0.0.2
pkgrel=1
pkgdesc="Terminal UI for browsing and selecting Nerd Font icons"
arch=('x86_64')
url="https://github.com/WillyV3/icon-picker"
license=('MIT')
depends=('wl-clipboard')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/WillyV3/$pkgname/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "$pkgname-$pkgver"
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export CGO_LDFLAGS="${LDFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
    go build -ldflags "-linkmode external -extldflags \"${LDFLAGS}\"" -o icon-picker .
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 icon-picker "$pkgdir/usr/bin/icon-picker"
    install -Dm755 icons "$pkgdir/usr/bin/icons"
    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"

    # Create cache directory structure
    install -dm755 "$pkgdir/usr/share/$pkgname"
}
