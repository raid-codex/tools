<?php
/**
* Plugin Name: Related posts
* Plugin URI: https://github.com/raid-codex/tools
* Description: Related posts
* Version: 1.0
* Author: Geoffrey Bauduin
* Author URI: https://github.com/geoffreybauduin
*/

function wpb_related_pages($atts) {
    $atts = shortcode_atts(
        array(
            'limit' => 2,
            'strategy' => "best-match",
            'shuffle' => true,
            'override_excerpt_meta' => null,
        ), $atts, "related-posts"
    );
    global $post;
    $orig_post = $post;
    $category_ids = wp_get_post_categories($post->ID);
    $html = "";
    if ($category_ids) {
        $args = array(
            'post_type' => 'page',
            'category__in' => $category_ids,
            'post__not_in' => array($post->ID),
            'nopaging' => true,
            'posts_per_page' => -1,
        );
        $query = new WP_Query($args);
        if ($query->have_posts())
        {
            $posts_common_categories = array();
            $atts["category-ids"] = $category_ids;
            $posts = related_posts_get_posts($query, $atts);
            $idx = 0;
            $html .= '<div class="related-posts"><div class="row">';
            foreach ($posts as $post)
            {
                if ($idx >= $atts['limit'])
                {
                    break;
                }
                $link = get_permalink($post);
                $title = get_the_title($post);
                $readMore = '<a href="'.$link.'" title="'.$title.'" rel="bookmark">Read more</a>';
                if ($atts['override_excerpt_meta'])
                {
                    $excerpt = get_post_meta($post->ID, $atts['override_excerpt_meta'], true);
                    $excerpt .= " ".$readMore;
                }
                else
                {
                    $excerpt = get_the_excerpt($post);
                    $excerpt = str_replace("[&hellip;]", '... '.$readMore, $excerpt);
                }
                $thumb = get_the_post_thumbnail($post, 'full');
                $html .= <<<EOF
<div class="col-xs-12">
    <div class="row" style="margin-top: 10px;">
        <div class="col-xs-12 col-lg-4 col-sm-6">
            <div class="related-posts-thumbnail">
                <a href="$link" title="$title" rel="bookmark">
                    $thumb
                </a>
            </div>
        </div>
        <div class="col-xs-12 col-lg-8 col-sm-6">
            <div class="related-posts-content">
                <a href="$link" title="$title" rel="bookmark">
                    <span class="h5">$title</span>
                </a>
                <p>$excerpt</p>
            </div>
        </div>
    </div>
</div>
EOF;
                $idx++;
            }
            $html .= '</div></div>';
        }
    }
    wp_reset_query();
    $post = $orig_post;
    return $html;
}

add_shortcode('related-posts', 'wpb_related_pages');

function related_posts_get_posts($query, $atts)
{
    if ($atts['strategy'] == "best-match")
    {
        return related_posts_get_posts_best_match($query, $atts);
    }
    return related_posts_get_posts_default($query, $atts);
}

function related_posts_get_posts_best_match($query, $atts)
{
    global $post;
    while($query->have_posts())
    {
        $query->the_post();
        $post_category_ids = wp_get_post_categories($post->ID);
        $intersect = array_intersect($post_category_ids, $atts["category-ids"]);
        $length = sizeof($intersect);
        if (!isset($posts_common_categories[$length]))
        {
            $posts_common_categories[$length] = array();
        }
        $posts_common_categories[$length][] = $post;
    }
    krsort($posts_common_categories);
    $posts = array();
    foreach ($posts_common_categories as $common => $posts_found)
    {
        if ($atts["shuffle"])
        {
            shuffle($posts_found);
        }
        $posts = array_merge($posts, $posts_found);
    }
    return $posts;
}

function related_posts_get_posts_default($query, $atts)
{
    global $post;
    $posts = array();
    while($query->have_posts())
    {
        $query->the_post();
        $posts[] = $post;
    }
    return $posts;
}

?>