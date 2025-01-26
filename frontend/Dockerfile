# Use a lightweight web server to serve static files
FROM nginx:alpine

# Copy the HTML file to the Nginx server
COPY index.html /usr/share/nginx/html/index.html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80
ENTRYPOINT [ "nginx" ]
CMD [ "-g", "daemon off;" ]
