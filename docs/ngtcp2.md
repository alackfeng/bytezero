


## ngtcp2.
`````

export NGTCP2_INSTALL=/home/vagrant/ngtcp2/install
cd /home/vagrant/ngtcp2/

git clone --depth 1 -b OpenSSL_1_1_1o+quic https://github.com/quictls/openssl
cd openssl
./config enable-tls1_3 --prefix=$NGTCP2_INSTALL/openssl
make -j$(nproc)
make install_sw

cd ..
git clone https://github.com/ngtcp2/nghttp3
cd nghttp3
autoreconf -i
./configure --prefix=$NGTCP2_INSTALL/nghttp3 --enable-lib-only
make -j$(nproc)
make install

cd ..
cd libev-4.33
./configure --prefix=$NGTCP2_INSTALL/libev
make -j$(nproc)
make install

cd ..
git clone https://github.com/ngtcp2/ngtcp2
cd ngtcp2
autoreconf -i
./configure CXX="/opt/rh/devtoolset-8/root/usr/bin/g++" CC="/opt/rh/devtoolset-8/root/usr/bin/gcc" PKG_CONFIG_PATH=$NGTCP2_INSTALL/openssl/lib/pkgconfig:$NGTCP2_INSTALL/nghttp3/lib/pkgconfig LDFLAGS="-Wl,-rpath,$NGTCP2_INSTALL/openssl/lib" LIBEV_CFLAGS="-I$NGTCP2_INSTALL/libev/include " LIBEV_LIBS="-L$NGTCP2_INSTALL/libev/lib -lev" --prefix=$NGTCP2_INSTALL/ngtcp2  --with-libev
make -j$(nproc) check
make install


`````
