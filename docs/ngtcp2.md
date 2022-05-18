


## ngtcp2.
`````

export NGTCP2_INSTALL=/home/vagrant/ngtcp2/install

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
make -j$(nproc) check
make install



cd ..
git clone https://github.com/ngtcp2/ngtcp2
cd ngtcp2
autoreconf -i
./configure PKG_CONFIG_PATH=$NGTCP2_INSTALL/openssl/lib/pkgconfig:$NGTCP2_INSTALL/nghttp3/lib/pkgconfig LDFLAGS="-Wl,-rpath,$NGTCP2_INSTALL/openssl/lib" --prefix=$NGTCP2_INSTALL/ngtcp2
make -j$(nproc) check


`````
