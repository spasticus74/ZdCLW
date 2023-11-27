<script>
    export let productOptions = [{"text":"Select a product..."},{"text":"API"},{"text":"Argus"},{"text":"Dispatch 5"},{"text":"Dispatch 6"},{"text":"Dispatch Underground"},{"text":"DTM"},{"text":"Frontrunner"},{"text":"Guided Spotting"},{"text":"KMAP"},{"text":"Minecare 2"},{"text":"Minecare 3"},{"text":"mRoc"},{"text":"Pegasys"},{"text":"Provision 2"},{"text":"Provision 3"},{"text":"Other Product"}];
    export let display_func = a => a;
    export let index = 0;
    export let value;
  
    $: {
      value = productOptions[index];
    }
</script>
  
<select bind:value={index}>
    {#each productOptions as productOption, i}
      <option value={i}>{display_func(productOption)}</option>
    {/each}
</select>
  