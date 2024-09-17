<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { GetCart, ProcessOrder, SetDishCount } from "../../client/cart";
  import { GetDishes } from "../../client/dish";
  import CartItem from "./cart_item.svelte";
  import {
    initMainButton,
    initBackButton,
    retrieveLaunchParams,
  } from "@telegram-apps/sdk";
  const { initData } = retrieveLaunchParams();
  import { navigate } from "svelte-routing";
  import { FormatPriceDefault } from "../../utils/format_price";

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
      let result = await ProcessOrder(initData.user.id, wishes);
      mainButton.enable();
      if (!result) {
        return;
      }
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

<h3>Корзина:</h3>
<div class="line"></div>
<div class="cart">
  {#each dishes as dish}
    <CartItem bind:dish />
  {/each}
  <div class="line"></div>
</div>
<div class="total_block">
  <div class="total_text">Итого:</div>
  <div class="spacer"></div>
  <div class="total_price">{FormatPriceDefault(total)}</div>
</div>
<div class="wishes">
  <textarea bind:value={wishes} class="wishes-input"></textarea>
</div>

<style>
  .cart {
    font-weight: normal;
    font-size: medium;
  }
  .wishes {
    width: 100vw;
    height: 30vh;
    text-align: center;
  }

  .wishes-input {
    background-color: var(--tg-theme-secondary-bg-color);
    color: var(--tg-theme-text-color);
    font-size: 14px;

    width: 80vw;
    height: 20vh;
    border: 0px;
    border-radius: 3%;
    resize: none;
    outline: none;
  }

  .line {
    border-radius: 30%;
    height: 5px;
    background-color: var(--tg-theme-secondary-bg-color);
  }

  .total_block {
    display: inline-flex;
    width: 100%;
  }

  .total_text {
    text-align: left;
    padding-left: var(--container-padding-lr);
  }
  .total_price {
    text-align: right;
    padding-right: var(--container-padding-lr);
  }
</style>
