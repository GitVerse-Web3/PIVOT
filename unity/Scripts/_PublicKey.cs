using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class _PublicKey : PublicKey
{
	public _PublicKey()
	{
		key = "key_" + GetHashCode().ToString();
	}

	public string key { get; set; }
}