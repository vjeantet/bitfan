## LRUCACHE

Lrucache is a powerful key/value store for Go.

You can use it almost the same way as an in-memory dictionary (`map[string]interface{}`) but there are many important differences in implementation.

## Features

<table>

<thead>
<tr>
<th></th>
<th>
<code>
map[string]interface{}
</code>
</th>
<th>lrucache</th>
</tr>
</thead>

<tbody>

<tr>
<th>thread-safe</th>
<td>no</td>
<td>yes</td>
</tr>

<tr>
<th>maximum size</th>
<td>no</td>
<td>yes</td>
</tr>

<tr>
<th>OnMiss handler</th>
<td>no</td>
<td>yes</td>
</tr>

</tbody>
</table>

* purges least recently used element when full
* elements can report their own size
* everything is cacheable (`interface{}`)
* is a front for your persistent storage (S3, disk, ...) by using OnMiss hooks

Examples and API are on godoc:

<http://godoc.org/github.com/hraban/lrucache>

The licensing terms are described in the file LICENSE.
