<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";
  import { initBackButton } from "@telegram-apps/sdk";

  import Dish from "./dish.svelte";
  import { LoadCart } from "../../client/cart";
  import DishesGrid from "../components/dishes_grid.svelte";
  import CategoriesList from "../components/categories_list.svelte";

  let mainButton = LoadCart();
  let removeListFn;
  let removeBackButtonListFn;
  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];

  let dishesGrid: DishesGrid;
  let categoriesList: CategoriesList;
  onMount(async () => {
    removeListFn = mainButton.on("click", () => {
      navigate("/cart");
      mainButton.hide();
    });
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
    removeListFn();
    removeBackButtonListFn();
  });
</script>

<main class="dish">
  <CategoriesList bind:this={categoriesList} />
  <horizontalSpacer class="primary-bg y-5" />
  <DishesGrid pageLimit={10} DishComponent={Dish} bind:this={dishesGrid} />
</main>
