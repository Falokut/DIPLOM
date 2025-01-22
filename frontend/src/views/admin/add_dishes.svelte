<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import DishPreview from "./dish_preview.svelte";
  import { GetAllDishesCategories } from "../../client/dishes_categories";
  import { AddDish } from "../../client/dish";
  import { ToBase64 } from "../../utils/base64";
  import { navigate } from "svelte-routing";

  import { initMainButton, initBackButton, number } from "@telegram-apps/sdk";
  import TextInput from "../components/text_input.svelte";
  import ImageInput from "../components/image_input.svelte";
  import MultiSelectInput from "../components/multi_select_input.svelte";

  let dish = {
    name: "",
    categories: [],
    url: "",
    price: "",
  };

  var selectedCategories = [];
  var image: File = null;

  var dishesCategoriesMap: Map<string, number> = new Map<string, number>();
  const mainButtonRes = initMainButton();
  var mainButton = mainButtonRes[0];

  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];

  var removeMainButtonListFn;
  var removeBackButtonListFn;
  onMount(() => {
    mainButton.setParams({
      text: "Добавить блюдо",
      isVisible: true,
    });

    mainButton.enable();
    removeMainButtonListFn = mainButton.on("click", addDish);

    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/admin", { replace: true });
    });
  });
  onDestroy(() => {
    mainButton.hide();
    removeMainButtonListFn();
    removeBackButtonListFn();
  });

  async function loadDishesCategories() {
    let categories = await GetAllDishesCategories();
    let dishCategories = [];
    categories.forEach((value) => {
      dishCategories.push(value.name);
      dishesCategoriesMap.set(value.name, value.id);
    });
    return dishCategories;
  }

  async function addDish() {
    mainButton.disable();
    let price = Math.ceil(Number(dish.price) * 100) / 100;
    if (price < 80) {
      return;
    }
    let req = {
      name: dish.name,
      categories: dish.categories,
      price: price * 100,
      image: null,
    };
    if (image != null && image.size > 0) {
      await ToBase64(image).then((data) => (req.image = data));
    }

    let ok = await AddDish(req);
    if (ok) {
      window.location.reload();
      return;
    }
    mainButton.enable();
  }
</script>

<main>
  <h3>Добавить блюдо</h3>
  <section class="add-dish-container">
    <TextInput bind:value={dish.name} label={"название:"} />
    <TextInput bind:value={dish.price} label={"цена:"} />
    <ImageInput
      bind:outputUrl={dish.url}
      label={"картинка:"}
      bind:file={image}
      uploadLabel={"выбрать файл"}
    />
    {#await loadDishesCategories() then dishCategories}
      <MultiSelectInput
        options={dishCategories}
        label={"категории:"}
        bind:selected={selectedCategories}
        onchange={() => {
          dish.categories = [];
          selectedCategories.forEach((name) => {
            let id = dishesCategoriesMap.get(name);
            dish.categories.push(id);
          });
        }}
      />
    {/await}
  </section>
  <horizontalSpacer class="primary-bg y-5" />
  <section class="dish-preview">
    {#if dish.url != "" && dish.name != "" && dish.price != ""}
      <h3>Предосмотр:</h3>
      <DishPreview bind:dish />
    {/if}
  </section>
</main>

<style>
  .add-dish-container {
    background-color: var(--secondary-bg-color);
    display: flex;
    flex-direction: column;
  }

  .dish-preview {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
  }
</style>
