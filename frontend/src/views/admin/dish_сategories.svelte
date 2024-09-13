<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { GetDishesCategories } from "../../client/dishes_categories";
  import { retrieveLaunchParams } from "@telegram-apps/sdk";
  import { navigate } from "svelte-routing";

  import { initBackButton } from "@telegram-apps/sdk";
  import DishCategory from "./dish_category.svelte";
  import AddDishCategory from "./add_dish_category.svelte";

  const { initData } = retrieveLaunchParams();
  const notAllowed = initData === undefined || initData.user === undefined;

  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];
  let categories = [];
  var removeBackButtonListFn;
  onMount(async () => {
    if (notAllowed) {
      window.close();
      return;
    }
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/admin", { replace: true });
    });

    categories = await GetDishesCategories();
  });
  onDestroy(() => {
    removeBackButtonListFn();
  });
  function remove(id: number) {
    categories = categories.filter((v) => v.id != id);
  }
</script>

<div class="dish_category_div">
  {#each categories as category}
    <DishCategory
      bind:category
      remove={() => {
        remove(category.id);
      }}
    ></DishCategory>
  {/each}
</div>

<div class="add_dish_category_container">
  <AddDishCategory
    OnAdd={() => {
      window.location.reload();
    }}
  ></AddDishCategory>
</div>

<style>
  .dish_category_div {
    display: grid;
    grid-auto-flow: row;
    gap: 2rem;
    width: 90vw;
    background-color: var(--tg-theme-secondary-bg-color);
    border-radius: 5%;
  }
  .add_dish_category_container {
    display: flex;

    padding-top: 1rem;
    width: 90vw;
    height: auto;
    min-height: 30vh;
    padding-bottom: 1rem;
  }
</style>
