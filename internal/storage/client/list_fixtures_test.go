package client_test

const listXMLPage1 = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
	<Name>my-bucket</Name>
	<Prefix>photos/</Prefix>
	<IsTruncated>true</IsTruncated>
	<NextContinuationToken>token-2</NextContinuationToken>
	<Contents>
		<Key>photos/1.jpg</Key>
		<Size>100</Size>
		<ETag>&quot;etag1&quot;</ETag>
	</Contents>
	<Contents>
		<Key>photos/2.jpg</Key>
		<Size>200</Size>
		<ETag>&quot;etag2&quot;</ETag>
	</Contents>
	<CommonPrefixes>
		<Prefix>photos/2025/</Prefix>
	</CommonPrefixes>
</ListBucketResult>`

const listXMLPage2 = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
	<Name>my-bucket</Name>
	<Prefix>photos/</Prefix>
	<IsTruncated>false</IsTruncated>
	<Contents>
		<Key>photos/3.jpg</Key>
		<Size>300</Size>
		<ETag>&quot;etag3&quot;</ETag>
	</Contents>
</ListBucketResult>`
