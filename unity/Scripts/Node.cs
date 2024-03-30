using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Zenject;

public class Node : MonoBehaviour, ICommit
{
	public float scaleFactor = 10;
	public float speed = 0.1f;
	public float r = 6;
	public float deltaY = 1;

	[Inject]
	ICommit _commit;

	[Inject]
	TextMesh _textMesh;

	public long modelHashID => _commit.modelHashID;

	public PublicKey author => _commit.author;

	public long authorSignature => _commit.authorSignature;

	public DateTime timestamp => _commit.timestamp;

	public string commitMessage => _commit.commitMessage;

	public ICommit parentModel
	{
		get => _commit.parentModel;
		set
		{

			_commit.parentModel = value;

		}
	}

	public double compressionRatio
	{
		get => _commit.compressionRatio;
		set
		{
			_commit.compressionRatio = value;
			updateScale();
		}
	}

	Tag ICommit.tag => _commit.tag;

	public bool checkValid()
	{
		return _commit.checkValid();
	}

	public void chosenToBeHead()
	{
		_commit.chosenToBeHead();
		var v = this.transform.position;
		v.x = 0;
		v.z = 0;
		this.transform.position = v;
	}

	public byte[] getFullModel()
	{
		return _commit.getFullModel();
	}

	public void rebaseToMaster(ICommit head)
	{

		Node h = (Node)head;
		updateY(h);
		_commit.rebaseToMaster(head);
	}

	void updateY(Node head)
	{
		var v = this.transform.position;
		float y = head.transform.localScale.y + deltaY + this.transform.localScale.y;
		v.y = y;
		this.transform.position = v;
	}

	// Start is called before the first frame update
	void Start()
	{
		updateDisplay();
	}

	// Update is called once per frame
	void Update()
	{
		var v = this.transform.position;
		Vector2 rr = new Vector2(v.x, v.z);
		var m = rr.magnitude;
		Vector3 target = new Vector3(v.x / m, v.y, v.z / m);
		v += (target - v) * speed;
		this.transform.position = v;

	}

	void updateScale()
	{
		this.transform.localScale = new Vector3(
			(float)compressionRatio * scaleFactor,
			(float)compressionRatio * scaleFactor,
			(float)compressionRatio * scaleFactor
			);
	}

	void updateDisplay()
	{
		_textMesh.text =
			"id: " + modelHashID
			+ "\n message: " + commitMessage
			+ "\n author: " + author
			+ "\n c: " + compressionRatio;
	}
}
