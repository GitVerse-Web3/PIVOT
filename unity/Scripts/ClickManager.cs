using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Zenject;

public class ClickManager : MonoBehaviour
{
	[Inject]
	BranchManager _branchManager;

	[Inject]
	public new Camera camera;


	RaycastHit hit;

	void rebase(Node node)
	{
		node.rebaseToMaster(_branchManager.masterHead);
	}



	void Start()
	{


	}

	// Update is called once per frame
	void Update()
	{
		if (Input.GetMouseButtonDown(0))
		{
			//Ray ray = camera.ScreenPointToRay(Input.mousePosition);
			Ray ray = camera.ScreenPointToRay(new Vector3(
				Screen.width / 2, Screen.height / 2, 0
				)
				);

			if (Physics.Raycast(ray, out hit))
			{
				Transform objectHit = hit.transform;

				Node node = objectHit.GetComponent<Node>();

				rebase(node);
			}
		}
	}
}
