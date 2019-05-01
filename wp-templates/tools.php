<?php

$starsPerGrade = array(
    "D" => array("stars" => 0, title => "not usable"),
    "C" => array("stars" => 1, title => "viable"),
    "B" => array("stars" => 2, title => "good"),
    "A" => array("stars" => 3, title => "exceptional"),
    "S" => array("stars" => 4, title => "top tier"),
    "SS" => array("stars" => 5, title => "god tier"),
);

function get_image_url_by_slug($slug, $size="thumbnail") {
    $args = array(
      'post_type' => 'attachment',
      'name' => sanitize_title($slug),
      'posts_per_page' => 1,
      'post_status' => 'inherit',
    );
    $_header = get_posts( $args );
    $header = $_header ? array_pop($_header) : null;
    return $header ? wp_get_attachment_image($header->ID, $size) : '';
}

function getStarsForGrade($grade) {
    return get_champion_grade($grade);
}

function get_champion_grade ( $grade ) {
    if ($grade == "")
    {
        return '<span class="champion-rating-none">No ranking yet</span>';
    }
    global $starsPerGrade;
    $str = '';
    for ($i = 0; $i < 5; $i++) {
        if ($i < $starsPerGrade[$grade]["stars"]) {
            $str = $str.'<i class="fas fa-star"></i>';
        }
        else {
            $str = $str.'<i class="far fa-star"></i>';
        }
    }
    return '<span class="champion-rating champion-rating-'.$grade.'" title="'.$starsPerGrade[$grade]["title"].'">'.$str.'</span>';
}

function get_permalink_by_slug( $slug ) {
    return get_page_link(get_page_by_path($slug));
}

function get_champion_rarity( $rarity ) {
    $classRarity = strtolower($rarity);
    return '<span class="champion-rarity champion-rarity-'.$classRarity.'">'.$rarity.'</span>';
}

?>