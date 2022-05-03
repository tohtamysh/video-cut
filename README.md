# video-cut

videocut --start=10 --stop=20 --file=path_to_video_file --output=path_to_output_folder

videocut -s 10 -e 20 -f path_to_video_file -o path_to_output_folder

## Docker example

docker run --rm -v $(pwd):/media ghcr.io/tohtamysh/video-cut /app/videocut -s 10 -e 120 -f /media/orig.mp4 -o /media
