# Markdown Blog

Markdown Blog is a simple and efficient application written in Go (Golang) that enables users to create and manage a blog website. Users can write their blog posts using Markdown, and the application handles rendering and serving the blog content by using HTML templates and static files.

## Features

- **Markdown Support**: Create and edit blog posts using Markdown for an intuitive writing experience.
- **Customizable Templates**: Modify or extend the look and feel of your blog using customizable HTML templates and CSS.
- **Tag and Category System**: Organize blog posts with tags and categories for better content management.

## Installation

### Prerequisites

- [Go](https://go.dev)
- [Node](https://nodejs.org)
- [TailwindCSS](https://tailwindcss.com)
- [DaisyUI](https://daisyui.com)

## Configuration

```
SUPERUSER_EMAIL=
CONTACT_EMAIL=

AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
BACKUP_ACCESS_KEY_ID=
BACKUP_SECRET_ACCESS_KEY=
CLOUDFRONT_DOMAIN=
APP_BUCKET=
BACKUP_BUCKET=
AWS_REGION=

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
```

## Customization

You can modify the templates in the `internal/templates/` directory to change the design of your blog. The following files are available:
- `pages/article.html`: Template for individual articles.
- `pages/index.html`: Template for the homepage.

Feel free to add your own CSS and JavaScript files in the `public/static/` directory.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- Markdown parsing powered by [Blackfriday](https://github.com/russross/blackfriday).
