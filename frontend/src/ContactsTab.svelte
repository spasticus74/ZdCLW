<script>
    import { SearchContactByName, SearchContactByOrg } from "../wailsjs/go/main/App.js";

    let searchTerm = "";
	let columns = ["Name", "Company", "Site", "Role", "Email", "Phone #1", "Phone #2"]
    let contacts = [];
    let resultsVisible = false

    async function searchName() {
        contacts = [];
        let APIResponse = SearchContactByName(searchTerm)
        let ResponseJson = JSON.parse(await APIResponse)
        let tempContacts = [];
        ResponseJson.forEach(contact => {
            let x = [];
            x[0] = contact.name
            x[1] = contact.company
            x[2] = contact.site
            x[3] = contact.role
            x[4] = '<a href="mailto:' + contact.email +'">' + contact.email + '</a>'
            x[5] = contact.phone1
            x[6] = contact.phone2
            tempContacts.push(x)
        });
        if (tempContacts.length > 0) {
            contacts = tempContacts
            resultsVisible = true
        } else {
            resultsVisible = false
            contacts = []
        }
    }

    async function searchOrgs() {
        contacts = [];
        let APIResponse = SearchContactByOrg(searchTerm)
        let ResponseJson = JSON.parse(await APIResponse)
        let tempContacts = [];
        ResponseJson.forEach(contact => {
            let x = [];
            x[0] = contact.name
            x[1] = contact.company
            x[2] = contact.site
            x[3] = contact.role
            x[4] = '<a href="mailto:' + contact.email +'">' + contact.email + '</a>'
            x[5] = contact.phone1
            x[6] = contact.phone2
            tempContacts.push(x)
        })
        if (tempContacts.length > 0) {
            contacts = tempContacts
            resultsVisible = true
        } else {
            resultsVisible = false
            contacts = []
        }
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
    
<button on:click={searchName}>
    Search by Name
</button>   
<button on:click={searchOrgs}>
    Search by Org
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
{:else }
    <p>No results</p>
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