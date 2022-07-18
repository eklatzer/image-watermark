# image-watermark

`image-watermark` is used to watermark images.

## Usage

Flags:
```text
  -height_percentage int
        Percentage of the height of the watermark (relative to the image it is placed on) (default 10)
  -input string
        Path to the folder containing the input files (must be .jpg) (default "./in")
  -offset_x int
        Distance of the watermark to the left side of the image
  -offset_y int
        Distance of the watermark to the bottom side of the image
  -output string
        Path for the images with watermark (default "./out")
  -output_sizes string
        List of sizes in which the output images are stored (width in pixels) separated by comma. Special value: source (default "source")
  -watermark string
        Path and file of the watermark (must be .png) (default "watermark.png")
```

## Example Usage
Folder structure:
```
.
├── image-watermark
├── in
│   ├── IMG_8265.JPG
│   ├── IMG_8297.JPG
│   ├── IMG_8345.JPG
│   ├── IMG_8405.JPG
│   ├── IMG_8416.JPG
│   ├── IMG_8629.JPG
│   └── IMG_8728.JPG
└── watermark.png
```

```
./image-watermark -offset_x=60 -offset_y=60 -height_percentage=10 -output_sizes="source,1920,3000"
```

Generates the new folder `out` with all watermarked images in the various width's:
```
.
├── image-watermark
├── in
│   ├── IMG_8265.JPG
│   ├── IMG_8297.JPG
│   ├── IMG_8345.JPG
│   ├── IMG_8405.JPG
│   ├── IMG_8416.JPG
│   ├── IMG_8629.JPG
│   └── IMG_8728.JPG
├── out
│   ├── 1920
│   │   ├── IMG_8265.JPG
│   │   ├── IMG_8297.JPG
│   │   ├── IMG_8345.JPG
│   │   ├── IMG_8405.JPG
│   │   ├── IMG_8416.JPG
│   │   ├── IMG_8629.JPG
│   │   └── IMG_8728.JPG
│   ├── 3000
│   │   ├── IMG_8265.JPG
│   │   ├── IMG_8297.JPG
│   │   ├── IMG_8345.JPG
│   │   ├── IMG_8405.JPG
│   │   ├── IMG_8416.JPG
│   │   ├── IMG_8629.JPG
│   │   └── IMG_8728.JPG
│   └── source
│       ├── IMG_8265.JPG
│       ├── IMG_8297.JPG
│       ├── IMG_8345.JPG
│       ├── IMG_8405.JPG
│       ├── IMG_8416.JPG
│       ├── IMG_8629.JPG
│       └── IMG_8728.JPG
└── watermark.png
```

Thanks to [esportfire.com](https://esportfire.com) for providing the example watermark image.