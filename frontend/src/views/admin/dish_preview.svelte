<script lang="ts">
  import PreviewImage from "../components/preview_image.svelte";
  import { FormatPriceDefault } from "../../utils/format_price";

  export let dish = {
    url: "",
    name: "",
    categories: [],
    price: 0,
  };
</script>

<section class="dish">
  <div class="dish-image">
    <PreviewImage bind:url={dish.url} bind:alt={dish.name} />
  </div>
  <div class="dish-caption">{dish.name}</div>
  <div class="dish-caption">{FormatPriceDefault(dish.price)}</div>
</section>

<style>
  .dish-caption {
    margin-top: auto;
  }
  .dish-image {
    width: 40vw;
    height: 40vw;
  }

  .dish {
    background-color: var(--secondary-bg-color);
    border-radius: 25%;
    text-align: center;
    display: flex;
    flex-direction: column;
    width: max-content;
    justify-content: center;
    align-items: center;

    padding: 30px;
  }
</style>
