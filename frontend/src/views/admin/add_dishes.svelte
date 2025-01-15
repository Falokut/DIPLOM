<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import DishPreview from "./dish_preview.svelte";
  import { GetAllDishesCategories } from "../../client/dishes_categories";
  import { AddDish } from "../../client/dish";
  import { ToBase64 } from "../../utils/base64";
  import { navigate } from "svelte-routing";

  import { initMainButton, initBackButton } from "@telegram-apps/sdk";
  import TextInput from "../components/text_input.svelte";
  import NumInput from "../components/num_input.svelte";
  import ImageInput from "../components/image_input.svelte";
  import MultiSelectInput from "../components/multi_select_input.svelte";
  import TextAreaInput from "../components/text_area_input.svelte";

  let dish = {
    name: "",
    categories: [],
    description: "",
    url: "",
    price: 0,
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
    let req = {
      name: dish.name,
      description: dish.description,
      categories: dish.categories,
      price: dish.price,
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

<h3>Добавить блюдо</h3>
<div class="add_dish_container">
  <TextInput bind:value={dish.name} label={"название:"} />
  <NumInput bind:value={dish.price} label={"цена:"} min={8000} max={1000000} />
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
  <TextAreaInput bind:value={dish.description} label={"описание"} />
</div>

{#if dish.url != "" && dish.name != "" && dish.price > 0}
  <h3>Предосмотр:</h3>
  <div class="dish_preview">
    <DishPreview bind:dish />
  </div>
{/if}

<style>
  .add_dish_container {
    background-color: var(--tg-theme-secondary-bg-color);
    display: grid;
    grid-template-rows: 1fr;
    grid-auto-flow: row;
    gap: 1rem;
    border-radius: 8px;
    padding: 20px 10px 20px 10px;
    width: 90vw;
  }

  .dish_preview {
    height: 50vh;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    padding-top: 20px;
  }
</style>
