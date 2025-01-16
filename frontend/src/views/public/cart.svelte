<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    ClearCart,
    ObjectFromCart,
    GetCart,
    SetDishCount,
  } from "../../client/cart";
  import { ProcessOrder } from "../../client/order";
  import { GetDishes } from "../../client/dish";
  import CartItem from "./cart_item.svelte";
  import { initMainButton, initBackButton } from "@telegram-apps/sdk";
  import { navigate } from "svelte-routing";
  import { FormatPriceDefault } from "../../utils/format_price";
  import TextAreaInput from "../components/text_area_input.svelte";

  const mainButtonRes = initMainButton();
  let mainButton = mainButtonRes[0];

  const backButtonRes = initBackButton();
  let backButton = backButtonRes[0];

  let dishes = [];
  let total = 0;
  let removeMainButtonListFn;
  let removeBackButtonListFn;
  let wishes = "";

  onMount(async () => {
    let cartItems = GetCart();
    let dishesIds = [];
    cartItems.forEach((v, k) => {
      if (v == undefined || k == undefined || v == 0) return;
      dishesIds.push(k);
    });

    dishes = await GetDishes(dishesIds);
    let exists = [];
    dishesIds.forEach((id) => {
      let found = dishesIds.findIndex((dish) => dish.id == id) == -1;
      if (!found) {
        SetDishCount(id, 0);
        return;
      }
      exists.push(id);
    });

    dishes.forEach((dish) => {
      let count = cartItems.get(dish.id.toString());
      if (count == 0) {
        return;
      }
      total += dish.price * count;
    });

    mainButton.setParams({
      text: "перейти к оплате",
      isVisible: true,
    });

    mainButton.enable();
    removeMainButtonListFn = mainButton.on("click", async () => {
      mainButton.disable();
      const items = ObjectFromCart(GetCart());
      let result = await ProcessOrder(items, wishes);
      mainButton.enable();
      if (!result) {
        return;
      }
      ClearCart();
      navigate("/", { replace: true });
    });

    backButton.show();
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/", { replace: true });
    });
  });

  onDestroy(() => {
    removeMainButtonListFn();
    removeBackButtonListFn();
    mainButton.hide();
  });
</script>

<div class="page-container">
  <h3>Корзина:</h3>
  <div class="line"></div>
  <div class="cart">
    {#each dishes as dish}
      <CartItem bind:dish />
    {/each}
    <div class="line"></div>
  </div>
  <div class="total-block">
    <div class="total-text">Итого:</div>
    <div class="total-price">{FormatPriceDefault(total)}</div>
  </div>
  <div class="spacer"></div>
  <div class="wishes">
    <TextAreaInput bind:value={wishes} label={"пожелания"} />
  </div>
</div>

<style>
  .cart {
    font-weight: normal;
    width: 100%;
  }

  .page-container {
    width: 90vw;
    border-radius: 5px;
  }

  .spacer {
    height: 5vh;
    background-color: var(--tg-theme-bg-color);
  }

  .wishes {
    width: 100%;
    border: 0px;
    resize: none;
    outline: none;
  }

  .line {
    height: 5px;
    background-color: var(--tg-theme-bg-color);
  }

  .total-block {
    display: flex;
    padding: 10px;
    justify-content: space-between;
    align-items: center;
    flex-direction: row;
  }

  .total-price {
    text-align: right;
  }
</style>
