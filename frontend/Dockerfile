# Gunakan Nginx Alpine sebagai base image (Sangat ringan)
FROM nginx:alpine

# (Opsional) Set working directory
WORKDIR /usr/share/nginx/html

# Hapus file default nginx
RUN rm -rf ./*

# Salin semua file dari folder frontend di Windows ke dalam container
# Docker otomatis menangani perbedaan path slash (\ vs /)
COPY . .

# Port standar untuk web
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]