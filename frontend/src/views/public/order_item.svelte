<script lang="ts">
  import { UserOrder } from "../../client/order";
  import { slide } from "svelte/transition";
  import { quintOut } from "svelte/easing";
  import { FormatPriceDefault } from "../../utils/format_price";
  import { onMount } from "svelte";

  export let order: UserOrder;
  let showOrderItems = false;
  let orderStatus = "";
  let orderDate = "";
  onMount(() => {
    let createdAt = new Date(order.createdAt);
    orderDate = createdAt.toISOString().slice(0, 10);
    orderDate += " "+createdAt.getHours() + ":" + createdAt.getMinutes();
    switch (order.status) {
      case "PROCESS":
        orderStatus = "в процессе";
        break;
      case "PAID":
        orderStatus = "оплачен";
        break;
      case "CANCELED":
        orderStatus = "отменён";
        break;
      case "SUCCESS":
        orderStatus = "выполнен";
        break;
      default:
        orderStatus = "";
        break;
    }
  });
</script>

<div class="order_box">
  <div class="order_info">
    <div class="order_text">
      Заказ от {orderDate}
    </div>
    <div class="spacer"></div>
    <button
      class="show_hide_order_items_button"
      on:click={() => {
        showOrderItems = !showOrderItems;
      }}>v</button
    >
  </div>
  <div>
    {#if showOrderItems}
      <div
        class="order_items"
        transition:slide={{
          delay: 250,
          duration: 300,
          easing: quintOut,
          axis: "y",
        }}
      >
        {#each order.items as item}
          <div class="user_order_item">
            <div class="name_count">{item.name} x {item.count}</div>
            <div class="spacer"></div>
            <div class="total_price">
              {FormatPriceDefault(item.totalPrice)}
            </div>
          </div>
        {/each}
        <div class="line"></div>
        <div class="order_total">
          <div class="user_order_total">Итого:</div>
          {#if orderStatus != ""}
            <div class="spacer"></div>
            <div class="order_text">{orderStatus}</div>
          {/if}
          <div class="spacer"></div>
          <div>{FormatPriceDefault(order.total)}</div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .order_box {
    padding-top: 3px;
    --order-item-border-color: var(--tg-theme-text-color) * 0.1;
  }
  .order_total {
    display: flex;
    flex-direction: row;
    margin: var(--container-padding-lr);
  }
  .order_items {
    background-color: var(--tg-theme-bg-color);
    border-radius: 8px;
    margin: var(--container-padding-lr);
    border: 2px solid var(--tg-theme-text-color);
    background-clip: padding-box;
  }
  .order_info {
    padding-top: 3px;
    display: flex;
    flex-direction: row;
    align-items: baseline;
    width: 100%;
    text-align: start;
    font-weight: normal;
  }
  .order_text {
    margin-left: var(--container-padding-lr);
  }
  .line {
    width: 100%;
    height: 1px;
    background-color: var(--tg-theme-text-color);
  }
  .show_hide_order_items_button {
    display: inline-flex;
    justify-content: center;
    background-color: transparent;
    border: 1px solid var(--tg-theme-text-color);
    margin-right: var(--container-padding-lr);
    padding-top: 0;
    padding-bottom: 0;
    font-weight: normal;
  }
  .show_hide_order_items_button:hover {
    color: green;
  }
  .user_order_item {
    display: inline-flex;
    flex-direction: row;
    flex-wrap: nowrap;

    width: 100%;
  }
  .user_order_item .name_count {
    padding-left: var(--container-padding-lr);
    text-align: left;
  }
  .user_order_item .total_price {
    padding-right: var(--container-padding-lr);
    text-align: right;
  }
</style>
