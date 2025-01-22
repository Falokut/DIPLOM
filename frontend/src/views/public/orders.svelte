<script lang="ts">
  import { GetUserOrders } from "../../client/order";
  import OrderItem from "./order_item.svelte";

  import { initBackButton } from "@telegram-apps/sdk";
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";

  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];

  let removeBackButtonListFn = () => {};
  onMount(() => {
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/", { replace: true });
    });
    fetchOrders();
    window.addEventListener("scroll", handleScroll);
  });
  onDestroy(() => {
    removeBackButtonListFn();
  });

  let orders = [];
  const pageLimit = 10;
  let currentOffset = 0;
  let loading = false;
  let allOrdersFetched = false;
  async function fetchOrders() {
    if (loading || allOrdersFetched) return;

    loading = true;
    let fetchedOrders = await await GetUserOrders(currentOffset, pageLimit);
    orders = [...orders, ...fetchedOrders];
    currentOffset = orders.length;
    loading = false;
    if (fetchedOrders.length < pageLimit) {
      allOrdersFetched = true;
    }
  }

  let container;
  function handleScroll(event) {
    if (!container) return;
    const { scrollTop, scrollHeight, clientHeight } = container;

    if (scrollTop + clientHeight >= scrollHeight - 10) {
      fetchOrders();
    }
  }
</script>

<main bind:this={container} class="orders-container">
  {#each orders as order}
    <OrderItem bind:order />
  {/each}
</main>

<style>
  .orders-container {
    height: 95vh;
    border-radius: 8px;
    background-color: var(--secondary-bg-color);
  }
</style>
