<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { GetAllDishesCategories } from "../../client/dishes_categories";
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

    categories = await GetAllDishesCategories();
  });
  onDestroy(() => {
    removeBackButtonListFn();
  });
  function remove(id: number) {
    categories = categories.filter((v) => v.id != id);
  }
</script>

{#key categories}
  <div class="dish-categories">
    {#each categories as category}
      <DishCategory
        {category}
        remove={() => {
          remove(category.id);
        }}
      ></DishCategory>
    {/each}
  </div>
{/key}

<div class="add-dish-categories">
  <AddDishCategory
    OnAdd={(category) => {
      categories.push(category);
      categories = categories;
    }}
  />
</div>

<style>
  .dish-categories {
    display: flex;
    flex-flow: column;
    width: 90vw;
    background-color: var(--tg-theme-secondary-bg-color);
    border-radius: 5px;
    padding: 10px;
  }

  .add-dish-categories {
    display: flex;

    padding: 10px;
    width: 90vw;
    height: auto;
    min-height: 30vh;
  }
</style>
