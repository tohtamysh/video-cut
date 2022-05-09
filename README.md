
# video-cut

videocut --start=10 --stop=20 --file=path_to_video_file --output=path_to_output_folder

videocut -s 10 -e 20 -f path_to_video_file -o path_to_output_folder

## Docker example

docker run --rm -v $(pwd):/media ghcr.io/tohtamysh/video-cut /app/videocut -s 10 -e 120 -f /media/orig.mp4 -o /media

### Batch

Что бы нарезать сразу много кусков с одного видео

docker run --rm -v $(pwd):/media ghcr.io/tohtamysh/video-cut /app/videocut -f /media/orig.mp4 -o /media -b /media/link.txt

```text
1-01,1-30,file_name
2-00,3-00
```
