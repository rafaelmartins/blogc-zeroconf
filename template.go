package main

const (
	mainTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ SITE_TITLE }}{% block entry %}{% if TITLE != SITE_TITLE %} | {{ TITLE }}{% endif %}{% endblock %}</title>
    <meta name="generator" content="blogc {{ BLOGC_VERSION }}" />
    <meta property="og:site_name" content="{{ SITE_TITLE }}" />
    {%- block entry %}
    {%- ifdef TITLE %}
    <meta property="og:title" content="{{ TITLE }}" />
    {%- endif %}
    {%- ifdef SITE_SUBTITLE %}
    <meta property="og:description" content="{{ SITE_SUBTITLE }}" />
    <meta name="description" content="{{ SITE_SUBTITLE }}" />
    {%- endif %}
    {%- endblock %}
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css@3.0.1/github-markdown.min.css" />
    <style type="text/css">
        .markdown-body {
            box-sizing: border-box;
            min-width: 200px;
            max-width: 980px;
            margin: 0 auto;
            padding: 45px;
        }
        @media (max-width: 767px) {
            .markdown-body {
                padding: 15px;
            }
        }
    </style>
  </head>
  <body>
    <div class="container-lg px-3 my-5 markdown-body">
      {%- ifdef SITE_TITLE %}
      <h1><a href="{{ BASE_URL }}/">{{ SITE_TITLE }}</a></h1>
      {%- endif %}
      {%- block entry %}
      {%- if TITLE != SITE_TITLE %}
      <h1>{{ TITLE }}</h1>
      {%- endif %}
      {{ CONTENT }}
      {%- endblock %}
      {%- block listing_entry %}
      {{ CONTENT }}
      <hr/>
      <h2>{{ POSTS_TITLE }}{% ifndef POSTS_TITLE %}Posts{% endif %}</h2>
      {%- endblock %}
      {%- block listing_once %}
      <ul>
      {%- endblock %}
      {%- block listing %}
        <li><a href="{{ BASE_URL }}/{{ POSTS_PREFIX }}{% ifndef POSTS_PREFIX %}post{% endif %}/{{ FILENAME }}/">{{ TITLE }}</a></li>
      {%- endblock %}
      {%- block listing_once %}
      </ul>
      {%- endblock %}
    </div>
  </body>
</html>`
)
