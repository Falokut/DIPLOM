<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { GetAllDishesCategories } from "../../client/dishes_categories";
  import { navigate } from "svelte-routing";

  import { initBackButton } from "@telegram-apps/sdk";
  import DishCategory from "./dish_category.svelte";
  import AddDishCategory from "./add_dish_category.svelte";

  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];
  let categories = [];
  var removeBackButtonListFn;
  onMount(async () => {
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

<main>
  <h3>Категории блюд</h3>
  {#key categories}
    <div class="dish-categories">
      {#each categories as category}
        <DishCategory {category} {remove} />
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
</main>

<style>
  .dish-categories {
    display: flex;
    flex-flow: column;
    width: 90vw;
    background-color: var(--secondary-bg-color);
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
