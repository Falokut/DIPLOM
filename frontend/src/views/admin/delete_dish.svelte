<script lang="ts">
  import { DeleteDish } from "../../client/dish";
  import PreviewImage from "../components/preview_image.svelte";
  import { Dish } from "../../client/dish";
  import { FormatPriceDefault } from "../../utils/format_price";

  export let dish: Dish;
  export let deletedFunc = function (id: any) {};

  async function deleteDish() {
    let confirmed = confirm("Удалить " + dish.name + "?");
    if (!confirmed) return;
    await DeleteDish(dish.id);
    deletedFunc(dish.id);
  }
</script>

<div class="dish">
  <div class="dish-img-count">
    <div class="dish-image">
      <PreviewImage bind:url={dish.url} bind:alt={dish.name} />
    </div>
  </div>
  <div class="dish-caption">{dish.name}</div>
  <div class="dish-caption">{FormatPriceDefault(dish.price)}</div>
  <div class="dish-count-buttons">
    <button
      class="dish-delete-button btn"
      on:click={() => {
        deleteDish();
      }}>удалить</button
    >
  </div>
</div>

<style>
  .dish-image {
    width: 150px;
    height: 150px;
  }

  .dish-caption {
    margin-top: auto;
    font-size: medium;
  }

  .dish {
    background-color: var(--tg-theme-secondary-bg-color);
    border-radius: 5%;
    padding: 10px;
    margin: 5px;
    text-align: center;
    display: inline-block;
  }

  .dish-delete-button {
    background-color: rgb(182, 19, 19);
  }
</style>
