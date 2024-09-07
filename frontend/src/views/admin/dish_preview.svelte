<script lang="ts">
  import PreviewImage from "../components/preview_image.svelte";

  export let dish = {
    url: "",
    name: "",
    categories: "",
    price: 0,
  };
  var previewSize = 90;
</script>

<div class="dish">
  <PreviewImage bind:url={dish.url} bind:size={previewSize} bind:alt={dish.name}
  ></PreviewImage>
  <div class="dish-caption">{dish.name}</div>
  <div class="dish-caption">{(dish.price / 100).toFixed(2) + "₽"}</div>
  <div class="dish-caption">{dish.categories}</div>
</div>

<style>
  :root {
    --dish-width: 110px;
    --dish-height: 200px;
  }

  .dish-caption {
    font-size: small;
    margin-top: auto;
  }

  .dish {
    text-align: center;
    display: inline-block;

    width: var(--dish-width);
    height: var(--dish-height);
  }
</style>
