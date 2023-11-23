<script>
    let searchTerm = "";
	let columns = ["Name", "Email", "Phone Number"]
    let contacts = [];
    let resultsVisible = false

    function displaySearch() {
        resultsVisible = true
        contacts =  [
            ["John", "john@example.com", "(353) 01 222 3333"],
            ["Sarah", "sarah@gmail.com", "(01) 22 888 4444"],
            ["Afshin", "afshin@mail.com", "(353) 22 87 8356"]
        ]
    }

    function clearResults() {
        resultsVisible = false
        contacts =  [];
    }

</script>
<h2>
    Enter a search term and click 'Search'
</h2>
    
<input type="searchVal" id="searchEntry" name="searchEntry" bind:value={searchTerm}><br>
    
<button on:click={displaySearch}>
    Search
</button>    
<button on:click={clearResults}>
    Clear
</button><br><br>

{#if resultsVisible}
<table id="results" align="center">
    <tr>
		{#each columns as column}
			<th>{column}</th>
		{/each}
	</tr>

    {#each contacts as row}
    <tr>
        {#each row as cell}
            <td contenteditable="true" bind:innerHTML={cell} />
        {/each}
    </tr>
{/each}
</table>
{/if}
    
<style>
    button {
        border-radius: .25rem;
        cursor: pointer;
        font-weight: bold;
        background-color: #fff;
        box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
        max-width: 250px;
    }
    button:active {
        background-color: #efefef;
        box-shadow: 0 0px 0px 0 rgba(0, 0, 0, 0.1), 0 0px 0px 0 rgba(0, 0, 0, 0.06);	
    }
</style>