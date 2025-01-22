<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";
  import { initBackButton } from "@telegram-apps/sdk";

  import DeleteDish from "./delete_dish.svelte";
  import CategoriesList from "../components/categories_list.svelte";
  import DishesGrid from "../components/dishes_grid.svelte";

  let removeBackButtonListFn;

  let dishesGrid: DishesGrid;
  let categoriesList: CategoriesList;
  onMount(() => {
    const backButtonRes = initBackButton();
    var backButton = backButtonRes[0];
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/", { replace: true });
    });
    categoriesList.$on(
      categoriesList.CategoryChangedEventType,
      (e: CustomEvent) => {
        dishesGrid.selectedCategoryUpdated(e.detail.selectedCategory);
      }
    );
  });

  onDestroy(() => {
    removeBackButtonListFn();
  });
</script>

<main class="dish">
  <CategoriesList bind:this={categoriesList} />
  <horizontalSpacer class="primary-bg y-5" />
  <DishesGrid
    pageLimit={10}
    DishComponent={DeleteDish}
    bind:this={dishesGrid}
  />
</main>
