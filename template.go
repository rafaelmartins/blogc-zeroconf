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
        .markdown-body footer {
            font-size: 12px
        }
        .markdown-body footer hr {
            height: 0.08em;
            margin: 24px 0 10px;
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
      {%- ifndef HIDE_FOOTER %}
      <footer>
        <hr />
        <p>
          Powered by: <a href="https://blogc.rgm.io/">blogc {{ BLOGC_VERSION }}</a> |
          Built {% ifdef BLOGC_SYSINFO_USERNAME %}by {{ BLOGC_SYSINFO_USERNAME }}{% endif %}
          {%- ifdef BLOGC_SYSINFO_HOSTNAME %}{% ifdef BLOGC_SYSINFO_USERNAME %}@
          {%- else %}at {% endif %}{{ BLOGC_SYSINFO_HOSTNAME }}
          {%- ifdef BLOGC_SYSINFO_INSIDE_DOCKER %} (docker){% endif %} {% endif %}
          {%- ifdef BLOGC_RUSAGE_CPU_TIME %}in {{ BLOGC_RUSAGE_CPU_TIME }}
          {%- ifdef BLOGC_SYSINFO_DATETIME %} ({{ BLOGC_SYSINFO_DATETIME }} GMT){% endif %}
          {%- ifdef BLOGC_RUSAGE_MEMORY %}, {% endif %}{% endif %}
          {%- ifdef BLOGC_RUSAGE_MEMORY %}using {{ BLOGC_RUSAGE_MEMORY }}{% endif %}.
        </p>
      </footer>
      {%- endif %}
    </div>
  </body>
</html>`

	atomTemplate = `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title type="text">{{ SITE_TITLE }}</title>
  <id>{{ BASE_DOMAIN }}{{ BASE_URL }}/atom{% ifdef FILTER_TAG %}/{{ FILTER_TAG }}{% endif %}/index.xml</id>
  <updated>{{ DATE_FIRST_FORMATTED }}</updated>
  <link href="{{ BASE_DOMAIN }}{{ BASE_URL }}/" />
  <link href="{{ BASE_DOMAIN }}{{ BASE_URL }}/atom{% ifdef FILTER_TAG %}/{{ FILTER_TAG }}{% endif %}/index.xml" rel="self" />
  <author>
    <name>{{ AUTHOR_NAME }}</name>
    {%- ifdef AUTHOR_EMAIL %}
    <email>{{ AUTHOR_EMAIL }}</email>
    {%- endif %}
  </author>
  {%- ifdef SITE_SUBTITLE %}
  <subtitle type="text">{{ SITE_SUBTITLE }}</subtitle>
  {%- endif %}
  {%- block listing %}
  <entry>
    <title type="text">{{ TITLE }}</title>
    <id>{{ BASE_DOMAIN }}{{ BASE_URL }}/{{ POSTS_PREFIX }}{% ifndef POSTS_PREFIX %}post{% endif %}/{{ FILENAME }}/index.html</id>
    <updated>{{ DATE_FORMATTED }}</updated>
    <published>{{ DATE_FORMATTED }}</published>
    <link href="{{ BASE_DOMAIN }}{{ BASE_URL }}/{{ POSTS_PREFIX }}{% ifndef POSTS_PREFIX %}post{% endif %}/{{ FILENAME }}/index.html" />
    <author>
      <name>{{ AUTHOR_NAME }}</name>
      {%- ifdef AUTHOR_EMAIL %}
      <email>{{ AUTHOR_EMAIL }}</email>
      {%- endif %}
    </author>
    <content type="html"><![CDATA[{{ CONTENT }}]]></content>
  </entry>
  {%- endblock %}
</feed>`
)
